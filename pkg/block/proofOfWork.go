package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/Cijin/gochain/pkg/utils"
)

const maxNounce = math.MaxInt64

// 256 cause we are using Sha256
const maxBits = 256

// Higher this number goes, the harder it will be to mine a block
const currentDifficulty = 24
const targetBits = maxBits - currentDifficulty

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// so resulting hash from mining has to be lower than target
	target = target.Lsh(target, targetBits)

	pow := ProofOfWork{
		block:  b,
		target: target,
	}

	return &pow
}

/*
 * Use headers to prepare data to be hashed, including nounce
 */
func (pow *ProofOfWork) prepareData(nounce int) []byte {
	b := pow.block
	data := bytes.Join([][]byte{
		b.HashTransactions(),
		b.PrevBlockHash,
		utils.ConvertToHex(b.Timestamp),
		utils.ConvertToHex(int64(nounce)),
	}, []byte{})

	return data
}

/*
* Run: Mining new Blocks.
 * What is mining?
 * It's the process of finding a hash that meets certain criteria, in this
 * case the criteria is that the hash should be smaller than target.
 *
*/
func (pow *ProofOfWork) Run() (int, []byte) {
	var hash [32]byte
	var hashInt big.Int
	var nounce int

	fmt.Println("Mining: ")
	for nounce < maxNounce {
		data := pow.prepareData(nounce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		// hash found
		if hashInt.Cmp(pow.target) == -1 {
			break
		}

		nounce++
	}
	fmt.Println()

	return nounce, hash[:]
}
