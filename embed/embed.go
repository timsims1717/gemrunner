package embed

import _ "embed"

// JiveTalking our true font file
//
//go:embed Jive_Talking.ttf
var JiveTalking []byte

// PuzzleShader our puzzle shader
//
//go:embed puzzle-shader.frag.glsl
var PuzzleShader string
