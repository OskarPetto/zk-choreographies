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
                6603137281342370245152077357736893208879593028427255370614241496850536950198
            ),
            uint256(
                928917802896087421311038053984171570794528664276848074966225084917712414047
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    5057314616560879242665713079858562778914247998990091193162201504223910664181
                ),
                uint256(
                    7261148373235841565466228814099404030647901750034362015317985271713891548767
                )
            ],
            [
                uint256(
                    7254420719305568482088728247205440402264845791891387727460905963397387407642
                ),
                uint256(
                    13819764893767814785637814956820128433705832811849173957272664419117630599072
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    17624653287347740333318913329301569094078428833631504462865530818959555340436
                ),
                uint256(
                    14662929739164875729335469971704538850464321686894611822433654134822921510643
                )
            ],
            [
                uint256(
                    8612140919648460443710731929090259110259864437567773535206688906926758183302
                ),
                uint256(
                    6771246259910446984932389520977939312392789835497037275253221274842398638541
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    15342868398016673370826062787452514004943907261489028980707427197794483157656
                ),
                uint256(
                    8232328439296367947273217362631292430716979098375996816389233485140448385196
                )
            ],
            [
                uint256(
                    17652025364006214181392915243339034583431697065606149911036377239821235798861
                ),
                uint256(
                    21363879073854992120188680392670624626743725872689870397012308681748593517425
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
            3681689532338548033602393526296756592560728939476569888833331026990651434985
        ); // vk.K[0].X
        vk_x.Y = uint256(
            3589690837978307462629658530385009086997339692718531452925793021211031353148
        ); // vk.K[0].Y
        mul_input[0] = uint256(
            19357180431342988470555915503574380578137759779511361910960634456179190694455
        ); // vk.K[1].X
        mul_input[1] = uint256(
            16182604524292307836026835500303016246316390663787026542618302890033927033973
        ); // vk.K[1].Y
        mul_input[2] = input[0];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[1] * input[0]
        mul_input[0] = uint256(
            8401461525651675592945209012916110431732856821686065466250025778043085383273
        ); // vk.K[2].X
        mul_input[1] = uint256(
            18676825606856939878489337413306280802899157359382544382095512460277961402871
        ); // vk.K[2].Y
        mul_input[2] = input[1];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[2] * input[1]
        mul_input[0] = uint256(
            1049574762636745974374305094548142199022037515810990162306433946716306476181
        ); // vk.K[3].X
        mul_input[1] = uint256(
            8384867663301490951833165647296351527161950626344435975997243304873859791574
        ); // vk.K[3].Y
        mul_input[2] = input[2];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[3] * input[2]
        mul_input[0] = uint256(
            21670032629998261526383340640433732203218443806275499517540579334192900702943
        ); // vk.K[4].X
        mul_input[1] = uint256(
            1293826273500490646345955245574149793053558582640849126373320063549614685547
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
