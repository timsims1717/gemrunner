package embed

import _ "embed"

// JiveTalking our true font file
//
//go:embed Jive_Talking.ttf
var JiveTalking []byte

// ColorShader our color shader (for UI elements)
//
//go:embed color-shader.frag.glsl
var ColorShader string

// PuzzleShader our puzzle view shader
//
//go:embed puzzle-shader.frag.glsl
var PuzzleShader string

// WorldShader our world view shader
//
//go:embed world-shader.frag.glsl
var WorldShader string
