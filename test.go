
type plugBoard struct {
	connections [26]int
}

type Reflector struct {
	wiring [26]int
}

type Rotor struct {
	wiring          [26]int
	notch           int
	currentPosition int
	inverseWiring   [26]int
}

type RotorSet struct {
	rotors []*Rotor
}

type InputRotor struct {
	wiring          [26]int
	currentPosition int
}

type EnigmaMachine struct {
	Plugboard     plugBoard
	Reflector     Reflector
	RotorSet      RotorSet
	InputRotor    InputRotor
	EncryptRotor1 Rotor
	EncryptRotor2 Rotor
	EncryptRotor3 Rotor
	DecryptRotor1 Rotor
	DecryptRotor2 Rotor
	DecryptRotor3 Rotor
}

// NewRotor creates a new Rotor instance
func NewRotor(wiring [26]int, notch int) *Rotor {
	r := &Rotor{
		wiring:          wiring,
		notch:           notch,
		inverseWiring:   [26]int{},
		currentPosition: 0,
	}

	// create inverse wiring
	for i := range wiring {
		r.inverseWiring[wiring[i]] = i
	}

	return r
}

// Rotate rotates the rotor to the next position and returns true if the notch is reached
func (r *Rotor) Rotate() bool {
	r.currentPosition = (r.currentPosition + 1) % 26
	return r.currentPosition == r.notch
}

// encipher enciphers the input character using the given rotor and returns the resulting character
func (r *Rotor) encipher(c int) int {
	// apply current position
	c = (c + r.currentPosition) % 26
	// apply wiring
	c = r.wiring[c]
	// apply inverse wiring
	c = (c - r.currentPosition + 26) % 26
	return c
}

// decipher deciphers the input character using the given rotor and returns the resulting character
func (r *Rotor) decipher(c int) int {
	// apply inverse wiring
	c = (c + r.currentPosition) % 26
	c = r.inverseWiring[c]
	// apply current position
	c = (c - r.currentPosition + 26) % 26
	return c
}

// NewRotorSet creates a new RotorSet instance with the given rotors
func NewRotorSet(rotors ...*Rotor) *RotorSet {
	return &RotorSet{rotors}
}

// Rotate rotates the rotor set
func (rs *RotorSet) Rotate() bool {
	rotateNext := rs.rotors[0].Rotate()
	for i := 1; i < len(rs.rotors) && rotateNext; i++ {
		rotateNext = rs.rotors[i].Rotate()
	}
	return rotateNext
}

// encipher enciphers the input character using the rotor set
func (rs *RotorSet) encipher(c int) int {
	for _, rotor := range rs.rotors {
		c = rotor.encipher(c)
	}
	return c
}

// decipher deciphers the input character using the rotor set
func (rs *RotorSet) decipher(c int) int {
	for i := len(rs.rotors) - 1; i >= 0; i-- {
		c = rs.rotors[i].decipher(c)
	}
	return c
}

// NewInputRotor creates a new InputRotor instance
func NewInputRotor(wiring [26]int) *InputRotor {
	return &InputRotor{wiring, 0}
}
