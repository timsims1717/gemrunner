package constants

// Shader Names
const (
	ShaderNone = iota
	ShaderWatery
	ShaderHeat
	ShaderEndOfList
)

const (
	DarknessDist  = float32(0.18)
	DarknessWidth = float32(0.85)
	DarknessGrad  = float32(0.05)
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
