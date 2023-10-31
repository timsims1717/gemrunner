#version 330 core

in vec2  vTexCoords;
out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

// custom uniforms
uniform float uRedPrimary;
uniform float uGreenPrimary;
uniform float uBluePrimary;

uniform float uRedSecondary;
uniform float uGreenSecondary;
uniform float uBlueSecondary;

// explosion texture
uniform sampler2D uExpTexture;

void main() {
    // Get our current screen coordinate
    vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

    // get our current coordinates' color
    vec4 col = texture(uTexture, t);

    // primary
    if (col.r == 1. && col.b == 1. && col.g == 0.) {
        col = vec4(uRedPrimary, uGreenPrimary, uBluePrimary, 1.);
    } else if (col.r == 0. && col.b == 1. && col.g == 1.) {
        col = vec4(uRedSecondary, uGreenSecondary, uBlueSecondary, 1.);
    }
    fragColor = col;
}
