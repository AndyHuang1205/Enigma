package main

import (
	"fmt"
	"strings"
)

type Machine struct {
	ALPHABET  string
	rotors    []*Rotor
	reflector *Reflector
}

func NewMachine(rotors []*Rotor, reflector *Reflector) *Machine {
	return &Machine{
		ALPHABET:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		rotors:    rotors,
		reflector: reflector,
	}
}

func (m *Machine) Encipher(text string) string {
	var result strings.Builder
	for _, x := range strings.ToUpper(text) {
		result.WriteString(string(m.encipherCharacter(x)))
	}
	return result.String()
}

func (m *Machine) Decipher(text string) string {
	for _, rotor := range m.rotors {
		rotor.Reset()
	}
	return m.Encipher(text)
}

func (m *Machine) encipherCharacter(x rune) rune {
	if strings.IndexRune(m.ALPHABET, x) == -1 {
		return x
	}
	if !strings.ContainsRune(m.ALPHABET, x) {
		return x
	}
	contactIndex := strings.IndexRune(m.ALPHABET, x)
	for _, rotor := range m.rotors {
		contactLetter := rune(rotor.alphabet[contactIndex])
		x = rune(rotor.Encipher(byte(contactLetter)))
		contactIndex = strings.IndexRune(rotor.alphabet, x)
	}
	contactLetter := rune(m.ALPHABET[contactIndex])
	x = m.reflector.Reflect(contactLetter)
	contactIndex = strings.IndexRune(m.ALPHABET, x)
	for i := len(m.rotors) - 1; i >= 0; i-- {
		rotor := m.rotors[i]
		contactLetter := rune(rotor.alphabet[contactIndex])
		x = rune(rotor.Decipher(byte(contactLetter)))
		contactIndex = strings.IndexRune(rotor.alphabet, x)
	}
	m.rotors[0].Rotate(1)
	for i := 1; i < len(m.rotors); i++ {
		rotor := m.rotors[i]
		turnFrequency := len(m.ALPHABET) * i
		if m.rotors[i-1].rotations%turnFrequency == 0 {
			rotor.Rotate(10)
		}
	}
	return rune(m.ALPHABET[contactIndex])
}

type Rotor struct {
	initialOffset int
	alphabet      string
	rotations     int
	forwardMap    map[byte]byte
	reverseMap    map[byte]byte
}

func NewRotor(mappings string, offset int) *Rotor {
	r := &Rotor{alphabet: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		initialOffset: offset}
	r.Reset()
	r.forwardMap = make(map[byte]byte)
	r.reverseMap = make(map[byte]byte)
	for i := 0; i < len(r.alphabet); i++ {
		r.forwardMap[r.alphabet[i]] = mappings[i]
		r.reverseMap[mappings[i]] = r.alphabet[i]
	}
	return r
}

func (r *Rotor) Reset() {
	r.alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r.Rotate(r.initialOffset)
	r.rotations = 1
}

func (r *Rotor) Rotate(offset int) {
	for i := 0; i < offset; i++ {
		r.alphabet = r.alphabet[1:] + string(r.alphabet[0])
	}
	r.rotations = offset
}

func (r *Rotor) Encipher(character byte) byte {
	return r.forwardMap[character]
}

func (r *Rotor) Decipher(character byte) byte {
	return r.reverseMap[character]
}

type Reflector struct {
	mappings map[rune]rune
}

const ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewReflector(mappings string) *Reflector {
	refl := &Reflector{
		mappings: make(map[rune]rune),
	}
	for i, c := range ALPHABET {
		refl.mappings[c] = rune(mappings[i])
	}

	return refl
}

func (r *Reflector) Reflect(char rune) rune {
	return r.mappings[char]
}

func main() {
	// create three rotors with different settings
	rotor1 := NewRotor("WNBPXCKJSFGUMLHYVQRTEIODAZ", 1)
	rotor2 := NewRotor("ZXCDNYIOBTQARSMHLPWFJEUGVK", 2)
	rotor3 := NewRotor("IKZSJVAUQLTOPYXBGRWNMFECDH", 3)
	rotor4 := NewRotor("TWUAVHPMZGFDXBIJNYKECLQSRO", 4)
	rotor5 := NewRotor("SPXYOFBMDLUJIQEAGTZRKWNCHV", 5)

	// create a reflector
	reflector := NewReflector("TJHVXQAFRYEZMUWSCLONDBPGKI")

	// create a machine with the three rotors and the reflector
	machine := NewMachine([]*Rotor{rotor1, rotor2, rotor3, rotor4, rotor5}, reflector)

	var text string
	fmt.Print("Enter your message: ")
	_, _ = fmt.Scanf("%s", &text)

	encrypted := machine.Encipher(text)
	fmt.Printf("encrypted text is %s", encrypted)
	fmt.Println()
	// decipher the encrypted message
	decrypted := machine.Decipher(encrypted)
	fmt.Printf("decrypted text is %s", decrypted)

}
