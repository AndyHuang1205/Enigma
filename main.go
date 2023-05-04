package main

import (
	"fmt"
	"strings"
)

type enigma struct {
	ALPHABET  string
	rotors    []*Rotor
	reflector *Reflector
	plugboard *Plugboard
}

func NewEnigma(rotors []*Rotor, reflector *Reflector, plugboard *Plugboard) *enigma {
	return &enigma{
		ALPHABET:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		plugboard: plugboard,
		rotors:    rotors,
		reflector: reflector,
	}
}

func (m *enigma) Encipher(text string) string {
	var result strings.Builder
	for _, x := range strings.ToUpper(text) {
		result.WriteString(string(m.encipherCharacter(x)))
	}
	return result.String()
}

func (m *enigma) Decipher(text string) string {
	for _, rotor := range m.rotors {
		rotor.Reset()
	}
	return m.Encipher(text)
}

func (m *enigma) encipherCharacter(x rune) rune {
	contactIndex := strings.IndexRune(m.ALPHABET, x)
	x = m.plugboard.Translate(rune(m.ALPHABET[contactIndex]))
	for _, rotor := range m.rotors {
		contactLetter := rotor.alphabet[contactIndex]
		x = rune(rotor.Encipher(byte(contactLetter)))
		contactIndex = strings.IndexRune(rotor.alphabet, x)
	}
	x = m.reflector.Reflect(rune(m.ALPHABET[contactIndex]))
	contactIndex = strings.IndexRune(m.ALPHABET, x)
	for i := len(m.rotors) - 1; i >= 0; i-- {
		rotor := m.rotors[i]
		contactLetter := rotor.alphabet[contactIndex]
		x = rune(rotor.Decipher(byte(contactLetter)))
		contactIndex = strings.IndexRune(rotor.alphabet, x)
	}
	x = m.plugboard.Translate(rune(m.ALPHABET[contactIndex]))
	m.rotors[0].Rotate(1)
	return rune(m.ALPHABET[contactIndex])
}

type Plugboard struct {
	mappings map[rune]rune
}

func NewPlugboard(mappingPairs []string) *Plugboard {
	pb := &Plugboard{
		mappings: make(map[rune]rune),
	}
	for _, pair := range mappingPairs {
		pb.mappings[rune(pair[0])] = rune(pair[1])
		pb.mappings[rune(pair[1])] = rune(pair[0])
	}
	return pb
}

func (pb *Plugboard) Translate(char rune) rune {
	if mapped, ok := pb.mappings[char]; ok {
		return mapped
	}
	return char
}

type Rotor struct {
	initialOffset int
	alphabet      string
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
}

func (r *Rotor) Rotate(offset int) {
	for i := 0; i < offset; i++ {
		r.alphabet = r.alphabet[1:] + string(r.alphabet[0])
	}
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
	reflector := NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT")
	plugboard := NewPlugboard([]string{"REMJGQDXFSUIVYOBKTLPACNWZH"})
	// create a machine with the three rotors and the reflector
	enigma := NewEnigma([]*Rotor{rotor1, rotor2, rotor3, rotor4, rotor5}, reflector, plugboard)

	var text string
	fmt.Print("Enter your message: ")
	_, _ = fmt.Scanf("%s", &text)

	encrypted := enigma.Encipher(text)
	fmt.Printf("encrypted text is %s", encrypted)
	fmt.Println()
	// decipher the encrypted message
	decrypted := enigma.Decipher(encrypted)
	fmt.Printf("decrypted text is %s", decrypted)

}
