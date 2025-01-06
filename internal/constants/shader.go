package constants

// Shader Names
const (
	ShaderNone = iota
	ShaderWatery
	ShaderHeat
	ShaderEndOfList
)

var (
	ShaderSpeeds = map[int]float32{
		ShaderWatery: 1.,
		ShaderHeat:   0.1,
	}

	ShaderXs = map[int]float32{
		ShaderWatery: 0.0007,
		ShaderHeat:   1.,
	}

	ShaderYs = map[int]float32{
		ShaderWatery: 0.0025,
		ShaderHeat:   1.25,
	}

	ShaderCustom = map[int]float32{
		ShaderWatery: 40.,
		ShaderHeat:   0.3,
	}
)

const (
	ParticleNone = iota
	ParticleDust
	ParticleEndOfList
)
