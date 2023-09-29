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

contract InstantiationVerifier {
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
                8465287537563014785583463023764568958751842362050334387301689160737629623397
            ),
            uint256(
                646124641151761918031802525385194759093891174899056450193831803994069754955
            )
        );
        vk.beta2 = Pairing.G2Point(
            [
                uint256(
                    3403163530182447719096620765742318756854980222972768665567346136753507279413
                ),
                uint256(
                    6699240369583604488083166160405094641207007532382031358651464353971819693362
                )
            ],
            [
                uint256(
                    9295164758918094287654186922088679938765228103592020617753707487467615401172
                ),
                uint256(
                    7382030499029781811173494110963957204377127316618083749774391822558205609971
                )
            ]
        );
        vk.gamma2 = Pairing.G2Point(
            [
                uint256(
                    21274995313566303813515679390045665790344314882396001104213143579050162590020
                ),
                uint256(
                    10856719320771483391077161775705821377953082832010942356051386021748353742179
                )
            ],
            [
                uint256(
                    7712550766609872139572173990494194767503274783952497971720609411181965943078
                ),
                uint256(
                    13116667048173396606528559745875382603780214903539823904191212717639770617736
                )
            ]
        );
        vk.delta2 = Pairing.G2Point(
            [
                uint256(
                    12629077115402688897708983850477275599598608678089795923676048104565950148634
                ),
                uint256(
                    10406204274361350737293772025361949625038930151216064138560133895929112049262
                )
            ],
            [
                uint256(
                    18981774951206802813951308248635626844420822768253726589086386462193184674899
                ),
                uint256(
                    6809884287053030941200233754323723925210660859419148870118018536063519448121
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
            6414588376385136492898561652524743610137432459262500877987468462591087144802
        ); // vk.K[0].X
        vk_x.Y = uint256(
            4139655507176118184020875740563747779592104564378894377005266841522133435961
        ); // vk.K[0].Y
        mul_input[0] = uint256(
            7655498933248654954460251542035620682983424663844141589011683178284388194472
        ); // vk.K[1].X
        mul_input[1] = uint256(
            4350841758185144873347524808940107446781761589960893172216851822026420430542
        ); // vk.K[1].Y
        mul_input[2] = input[0];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[1] * input[0]
        mul_input[0] = uint256(
            3918225364725647533768595703198554473811228124607583281571882034015853070740
        ); // vk.K[2].X
        mul_input[1] = uint256(
            20015641521923439796156639794730331137653255016603988976775236608090765314011
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
