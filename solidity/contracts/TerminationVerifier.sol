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

contract TerminationVerifier {
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
                16801867173666121034812193937395916981489291810204506257245512319774998573761
            ),
            uint256(
                8568745283113500864753198567277937484461745327061273385802883572645261041211
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    10999073524810989769207731660862382183210195139299738625320987360721091068290
                ),
                uint256(
                    4783204641107145644901414549949476759134046540830795050086846122485733431511
                )
            ],
            [
                uint256(
                    15107997664178432276950689471774705251546250478014945340537264948409926018593
                ),
                uint256(
                    21132019489050373730564653898114024564277584194878117632625761070449444090505
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    16975581726320958261613365923899336762622041514966035635734990829103990152783
                ),
                uint256(
                    12931808893192625424688168113458417384412689747662029653431203503338612547187
                )
            ],
            [
                uint256(
                    7583969531433127846346747122058187269625205500677953085115176011999024030261
                ),
                uint256(
                    10490281365354050267661583789104636612866937111732976057577910886117389622556
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    2162355983405531647537388645718954718676332792572742999713365616852460305808
                ),
                uint256(
                    3611241607479416864761942273789070414127185742924757507392118117444333819337
                )
            ],
            [
                uint256(
                    21246422336416151721241773991145012503135521740915169790733303238204729048449
                ),
                uint256(
                    18706162942939007093962859655733201034698019509207715769146917013392772958849
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
        uint256[2] calldata input
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
            9431161144816782649202163535841796459211977716111919442818282189352976438303
        ); // vk.K[0].X
        vk_x.Y = uint256(
            17671951472279846396514325801421725517963703368067388572497045189789855146974
        ); // vk.K[0].Y
        mul_input[0] = uint256(
            8035992916796374212149624973370437975164485768078392970615213734860486982486
        ); // vk.K[1].X
        mul_input[1] = uint256(
            19965648399164453595189644479532901926428168067038935238618390709810523524273
        ); // vk.K[1].Y
        mul_input[2] = input[0];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[1] * input[0]
        mul_input[0] = uint256(
            2647796612697982729176819946854340298503178058005997526524306665570104083154
        ); // vk.K[2].X
        mul_input[1] = uint256(
            8684790097121100928736879505823351063268203460622006247698537910366205481916
        ); // vk.K[2].Y
        mul_input[2] = input[1];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[2] * input[1]

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
