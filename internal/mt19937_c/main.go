// Package mt19937_c is the C reference implementation of mt19937-64.
package mt19937_c

// // initializes mt[NN] with a seed
// void init_genrand64(unsigned long long seed);
//
// // initialize by an array with array-length
// // init_key is the array for initializing keys
// // key_length is its length
// void init_by_array64(unsigned long long init_key[],
//                      unsigned long long key_length);
//
// // generates a random number on [0, 2^64-1]-interval
// unsigned long long genrand64_int64(void);
//
// // generates a random number on [0, 2^63-1]-interval
// long long genrand64_int63(void);
//
// void benchmark_rand(unsigned long long n)
// {
//     init_genrand64(123456789LLU);
//     for (unsigned long long i = 0; i < n; i++)
//     {
//         genrand64_int64();
//     }
// }
//
// void benchmark_seed(unsigned long long n)
// {
//     for (unsigned long long i = 0; i < n; i++)
//     {
//         init_genrand64(123456789LLU);
//     }
// }
//
// void benchmark_seed_multi(unsigned long long n)
// {
//     static unsigned long long seed[8] =
//     {
//         0x1234567890123456LLU,
//         0x7890123456789012LLU,
//         0x3456789012345678LLU,
//         0x9012345678901234LLU,
//         0x5678901234567890LLU,
//         0x1234567890123456LLU,
//         0x7890123456789012LLU,
//         0x3456789012345678LLU,
//     };
//     for (unsigned long long i = 0; i < n; i++)
//     {
//         init_by_array64(seed, sizeof(seed) / sizeof(seed[0]));
//     }
// }
import "C"
import "testing"

// initializes mt[NN] with a seed
func InitGenRand64(seed uint64) {
	C.init_genrand64(C.ulonglong(seed))
}

// initialize by an array with array-length
// init_key is the array for initializing keys
// key_length is its length
func InitByArray64(key ...uint64) {
	C.init_by_array64((*C.ulonglong)(&key[0]), C.ulonglong(len(key)))
}

// generates a random number on [0, 2^64-1]-interval
func GenRand64Int64() uint64 {
	return uint64(C.genrand64_int64())
}

// generates a random number on [0, 2^63-1]-interval
func GenRand64Int63() int64 {
	return int64(C.genrand64_int63())
}

func BenchmarkRand(b *testing.B) {
	C.benchmark_rand(C.ulonglong(b.N))
}

func BenchmarkSeed(b *testing.B) {
	C.benchmark_seed(C.ulonglong(b.N))
}

func BenchmarkSeedMulti(b *testing.B) {
	C.benchmark_seed_multi(C.ulonglong(b.N))
}
