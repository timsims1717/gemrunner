#version 330 core

in vec2  vTexCoords;
in vec4  vColor;
out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

// color uniforms
uniform float uRedPrimary;
uniform float uGreenPrimary;
uniform float uBluePrimary;

uniform float uRedSecondary;
uniform float uGreenSecondary;
uniform float uBlueSecondary;

uniform float uRedDoodad;
uniform float uGreenDoodad;
uniform float uBlueDoodad;

uniform float uRedLiquidPrimary;
uniform float uGreenLiquidPrimary;
uniform float uBlueLiquidPrimary;

uniform float uRedLiquidSecondary;
uniform float uGreenLiquidSecondary;
uniform float uBlueLiquidSecondary;

// world uniforms
uniform int uMode;
uniform float uSpeed;
uniform float uTime;

void main() {
    // Get our current screen coordinate
    vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

    // get our current coordinates' color
    vec4 col = texture(uTexture, t);

    // primary, secondary, tertiary
    if (col.r == 1. && col.b == 1. && col.g == 0.) {
        col = vec4(uRedPrimary, uGreenPrimary, uBluePrimary, col.a);
    } else if (col.r == 0. && col.b == 1. && col.g == 1.) {
        col = vec4(uRedSecondary, uGreenSecondary, uBlueSecondary, col.a);
    } else if (col.r == 1. && col.b == 0. && col.g == 1.) {
        col = vec4(uRedDoodad, uGreenDoodad, uBlueDoodad, col.a);
    } else if (col.r == 0. && col.b == 1. && col.g == 0.) {
        col = vec4(uRedLiquidPrimary, uGreenLiquidPrimary, uBlueLiquidPrimary, col.a);
    } else if (col.r == 1. && col.b == 0. && col.g == 0.) {
        col = vec4(uRedLiquidSecondary, uGreenLiquidSecondary, uBlueLiquidSecondary, col.a);
    }
    fragColor = col;
}
