
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

library Pairing {

    uint256 constant PRIME_Q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

    struct G1Point {
        uint256 X;
        uint256 Y;
    }

    // Encoding of field elements is: X[0] * z + X[1]
    struct G2Point {
        uint256[2] X;
        uint256[2] Y;
    }

    /*
     * @return The negation of p, i.e. p.plus(p.negate()) should be zero.
     */
    function negate(G1Point memory p) internal pure returns (G1Point memory) {

        // The prime q in the base field F_q for G1
        if (p.X == 0 && p.Y == 0) {
            return G1Point(0, 0);
        } else {
            return G1Point(p.X, PRIME_Q - (p.Y % PRIME_Q));
        }
    }

    /*
     * @return The sum of two points of G1
     */
    function plus(
        G1Point memory p1,
        G1Point memory p2
    ) internal view returns (G1Point memory r) {

        uint256[4] memory input;
        input[0] = p1.X;
        input[1] = p1.Y;
        input[2] = p2.X;
        input[3] = p2.Y;
        bool success;

        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 6, input, 0xc0, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success case 0 { invalid() }
        }

        require(success,"pairing-add-failed");
    }


    /*
     * Same as plus but accepts raw input instead of struct
     * @return The sum of two points of G1, one is represented as array
     */
    function plus_raw(uint256[4] memory input, G1Point memory r) internal view {
        bool success;

        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 6, input, 0xc0, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success case 0 {invalid()}
        }

        require(success, "pairing-add-failed");
    }

    /*
     * @return The product of a point on G1 and a scalar, i.e.
     *         p == p.scalar_mul(1) and p.plus(p) == p.scalar_mul(2) for all
     *         points p.
     */
    function scalar_mul(G1Point memory p, uint256 s) internal view returns (G1Point memory r) {

        uint256[3] memory input;
        input[0] = p.X;
        input[1] = p.Y;
        input[2] = s;
        bool success;
        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 7, input, 0x80, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success case 0 { invalid() }
        }
        require (success,"pairing-mul-failed");
    }


    /*
     * Same as scalar_mul but accepts raw input instead of struct,
     * Which avoid extra allocation. provided input can be allocated outside and re-used multiple times
     */
    function scalar_mul_raw(uint256[3] memory input, G1Point memory r) internal view {
        bool success;

        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 7, input, 0x80, r, 0x60)
            // Use "invalid" to make gas estimation work
            switch success case 0 {invalid()}
        }
        require(success, "pairing-mul-failed");
    }

    /* @return The result of computing the pairing check
     *         e(p1[0], p2[0]) *  .... * e(p1[n], p2[n]) == 1
     *         For example,
     *         pairing([P1(), P1().negate()], [P2(), P2()]) should return true.
     */
    function pairing(
        G1Point memory a1,
        G2Point memory a2,
        G1Point memory b1,
        G2Point memory b2,
        G1Point memory c1,
        G2Point memory c2,
        G1Point memory d1,
        G2Point memory d2
    ) internal view returns (bool) {

        G1Point[4] memory p1 = [a1, b1, c1, d1];
        G2Point[4] memory p2 = [a2, b2, c2, d2];
        uint256 inputSize = 24;
        uint256[] memory input = new uint256[](inputSize);

        for (uint256 i = 0; i < 4; i++) {
            uint256 j = i * 6;
            input[j + 0] = p1[i].X;
            input[j + 1] = p1[i].Y;
            input[j + 2] = p2[i].X[0];
            input[j + 3] = p2[i].X[1];
            input[j + 4] = p2[i].Y[0];
            input[j + 5] = p2[i].Y[1];
        }

        uint256[1] memory out;
        bool success;

        // solium-disable-next-line security/no-inline-assembly
        assembly {
            success := staticcall(sub(gas(), 2000), 8, add(input, 0x20), mul(inputSize, 0x20), out, 0x20)
            // Use "invalid" to make gas estimation work
            switch success case 0 { invalid() }
        }

        require(success,"pairing-opcode-failed");

        return out[0] != 0;
    }
}

contract Verifier {

    using Pairing for *;

    uint256 constant SNARK_SCALAR_FIELD = 21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant PRIME_Q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

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
        vk.alfa1 = Pairing.G1Point(uint256(11272699396656621149232229614088905745788345103537525813306460323054264646778), uint256(16817996263099757740537057652046463920051636461558508264937765493902704032693));
        vk.beta2 = Pairing.G2Point([uint256(234545153119156913706859158321032542416676892723908747595984322067681462784), uint256(1831467819239774819223883733395234759670039766495509678654660754716758697583)], [uint256(6546777898400028835933375252023302532301805591750023233050246856531387256250), uint256(7880057529177785919094771809136684389954695752956874962358803534355819790652)]);
        vk.gamma2 = Pairing.G2Point([uint256(16385821886399323736150824623121450865410423727180906236913250758055419894831), uint256(19207541448395197061850310947062544606768538793647429985246454269946378441589)], [uint256(13169204806323090603300963007746938631026570641853304289631281042583365963561), uint256(6176881884095324508006771519432547460123391601191553076589935550099594935)]);
        vk.delta2 = Pairing.G2Point([uint256(17985970887126141473810056936219121642535536283557042314564786345910742631555), uint256(4773650611112981352514064228415672638394148026985470547300947073385106495004)], [uint256(1633770085707036097988730277587324121278237299153356301856969086806342305770), uint256(12005844376384245297681919678948916876363681842762370853603092335414493426393)]);
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
            require(input[i] < SNARK_SCALAR_FIELD,"verifier-gte-snark-scalar-field");
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

        vk_x.X = uint256(18580594726508932553137905811546379006398136315493197644371417035190828853384); // vk.K[0].X
        vk_x.Y = uint256(19128444718099263445904065908236620115986330439015846286307476582846610785734); // vk.K[0].Y
        mul_input[0] = uint256(12104254595372516879867137191355762103150919854879716576703721570498660740941); // vk.K[1].X
        mul_input[1] = uint256(16612528882585380793473180103216326406210983404567740603613030054055771195223); // vk.K[1].Y
        mul_input[2] = input[0];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[1] * input[0]
        mul_input[0] = uint256(16351196693638992149585321368060221512853604501890468093840902909868642646335); // vk.K[2].X
        mul_input[1] = uint256(3747350425940075484109396655640655724811164888261394514222118606692150833561); // vk.K[2].Y
        mul_input[2] = input[1];
        accumulate(mul_input, q, add_input, vk_x); // vk_x += vk.K[2] * input[1]

        return Pairing.pairing(
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
