// SPDX-License-Identifier: AML
//
// Copyright 2017 Christian Reitwiessner
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// 2019 OKIMS

pragma solidity ^0.8.0;

import "./Pairing.sol";

contract TransitionVerifier {
    using Pairing for *;

    uint256 constant SNARK_SCALAR_FIELD =
        21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant PRIME_Q =
        21888242871839275222246405745257275088696311157297823662689037894645226208583;

    struct VerifyingKey {
        Pairing.G1Point alfa1;
        Pairing.G2Point beta2;
        Pairing.G2Point gamma2;
        Pairing.G2Point delta2;
        // []G1Point IC (K in gnark) appears directly in verifyProof
    }

    struct Proof {
        Pairing.G1Point A;
        Pairing.G2Point B;
        Pairing.G1Point C;
    }

    function verifyingKey() internal pure returns (VerifyingKey memory vk) {
        vk.alfa1 = Pairing.G1Point(
            uint256(
                16561205661185773988728068589565671245176452287306071214269889069516199275224
            ),
            uint256(
                7836552449832406717809906350548491061731231826060207667150091583190013357414
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    15669337191749468472107281404660982690648177796924660615377601879121744604056
                ),
                uint256(
                    6383688092391356170407277366701932686370655788998383044617647036799202774635
                )
            ],
            [
                uint256(
                    21096331034327272829095255853393101562514044881697900528192308209201588278691
                ),
                uint256(
                    11962828491854380267981675472284191149185821241175929707017934996549834780716
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    9010721528389790823974416429640526391293745277790554063778386437816912602657
                ),
                uint256(
                    16981584776567153902049213328446865857603268909417418540906675242085394133809
                )
            ],
            [
                uint256(
                    18564891735674485822843853987392376777116342398210461586536087929920789097483
                ),
                uint256(
                    5846829421294473555140486978556583858891484379018664163983059735669635290114
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    4490308926372578513177196274398076381705027050547674663645752184464816183443
                ),
                uint256(
                    14024702071870724417131946063298108423347633914538118719872241819393958016634
                )
            ],
            [
                uint256(
                    5147879172127114009446359914811168066157732578498319712965645685002384238623
                ),
                uint256(
                    5648444900276935568366341872261806824578139935147995939964888231585874171155
                )
            ]
        );
    }

    // accumulate scalarMul(mul_input) into q
    // that is computes sets q = (mul_input[0:2] * mul_input[3]) + q
    function accumulate(
        uint256[3] memory mul_input,
        Pairing.G1Point memory p,
        uint256[4] memory buffer,
        Pairing.G1Point memory q
    ) internal view {
        // computes p = mul_input[0:2] * mul_input[3]
        Pairing.scalar_mul_raw(mul_input, p);

        // point addition inputs
        buffer[0] = q.X;
        buffer[1] = q.Y;
        buffer[2] = p.X;
        buffer[3] = p.Y;

        // q = p + q
        Pairing.plus_raw(buffer, q);
    }

    /*
     * @returns Whether the proof is valid given the hardcoded verifying key
     *          above and the public inputs
     */
    function verifyProof(
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[4] calldata input
    ) public view returns (bool r) {
        Proof memory proof;
        proof.A = Pairing.G1Point(a[0], a[1]);
        proof.B = Pairing.G2Point([b[0][0], b[0][1]], [b[1][0], b[1][1]]);
        proof.C = Pairing.G1Point(c[0], c[1]);

        // Make sure that proof.A, B, and C are each less than the prime q
        require(proof.A.X < PRIME_Q, "verifier-aX-gte-prime-q");
        require(proof.A.Y < PRIME_Q, "verifier-aY-gte-prime-q");

        require(proof.B.X[0] < PRIME_Q, "verifier-bX0-gte-prime-q");
        require(proof.B.Y[0] < PRIME_Q, "verifier-bY0-gte-prime-q");

        require(proof.B.X[1] < PRIME_Q, "verifier-bX1-gte-prime-q");
        require(proof.B.Y[1] < PRIME_Q, "verifier-bY1-gte-prime-q");

        require(proof.C.X < PRIME_Q, "verifier-cX-gte-prime-q");
        require(proof.C.Y < PRIME_Q, "verifier-cY-gte-prime-q");

        // Make sure that every input is less than the snark scalar field
        for (uint256 i = 0; i < input.length; i++) {
            require(
                input[i] < SNARK_SCALAR_FIELD,
                "verifier-gte-snark-scalar-field"
            );
        }

        VerifyingKey memory vk = verifyingKey();

        // Compute the linear combination vk_x
        Pairing.G1Point memory vk_x = Pairing.G1Point(0, 0);

        // Buffer reused for addition p1 + p2 to avoid memory allocations
        // [0:2] -> p1.X, p1.Y ; [2:4] -> p2.X, p2.Y
        uint256[4] memory add_input;

        // Buffer reused for multiplication p1 * s
        // [0:2] -> p1.X, p1.Y ; [3] -> s
        uint256[3] memory mul_input;

        // temporary point to avoid extra allocations in accumulate
        Pairing.G1Point memory q = Pairing.G1Point(0, 0);

        vk_x.X = uint256(
            2629782331778877778614732383302198416338952297054663449574339459244002316508
        ); // vk.K[0].X
        vk_x.Y = uint256(
            19843951072959953422048731978383242283260672671710618624488028892787825495860
        ); // vk.K[0].Y
        mul_input[0] = uint256(
            21814750932109174379714703241156630224276334630130401957759509550234294757337
        ); // vk.K[1].X
        mul_input[1] = uint256(
            5368732045901069056686219896235278026087024529495445647458746378629448165492
        ); // vk.K[1].Y
        mul_input[2] = input[0];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[1] * input[0]
        mul_input[0] = uint256(
            17746161279938112608069584574086379655183990069869622096481712890849809370421
        ); // vk.K[2].X
        mul_input[1] = uint256(
            12527782403123888345961813162670307012323638569069138918777994250644869429330
        ); // vk.K[2].Y
        mul_input[2] = input[1];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[2] * input[1]
        mul_input[0] = uint256(
            1589302332935460332014373547715392460079557875672746795006825665845168808638
        ); // vk.K[3].X
        mul_input[1] = uint256(
            20257922194603055952753889537289129666358247603418865235869522324229129117057
        ); // vk.K[3].Y
        mul_input[2] = input[2];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[3] * input[2]
        mul_input[0] = uint256(
            21063573067252683536391645832570616450915495206604940108389742572961138411052
        ); // vk.K[4].X
        mul_input[1] = uint256(
            7387670020394089743320148055389240117739911439053601906874735469904928541411
        ); // vk.K[4].Y
        mul_input[2] = input[3];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[4] * input[3]

        return
            Pairing.pairing(
                Pairing.negate(proof.A),
                proof.B,
                vk.alfa1,
                vk.beta2,
                vk_x,
                vk.gamma2,
                proof.C,
                vk.delta2
            );
    }
}
