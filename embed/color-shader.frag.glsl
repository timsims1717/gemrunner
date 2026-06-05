#version 330 core

in vec2  vTexCoords;
in vec4  vColor;
out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

// color uniforms
uniform vec3 uPrimary;
uniform vec3 uSecondary;
uniform vec3 uDoodad;
uniform vec3 uGoop;
uniform vec3 uLiquidPrimary;
uniform vec3 uLiquidSecondary;

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
    if (col.r == 1. && col.g == 0. && col.b == 1.) {
        col = vec4(uPrimary.r, uPrimary.g, uPrimary.b, col.a);
    } else if (col.r == 0. && col.g == 1. && col.b == 1.) {
        col = vec4(uSecondary.r, uSecondary.g, uSecondary.b, col.a);
    } else if (col.r == 1. && col.g == 1. && col.b == 0.) {
        col = vec4(uDoodad.r, uDoodad.g, uDoodad.b, col.a);
    } else if (col.r == 0. && col.g == 0. && col.b == 1.) {
        col = vec4(uLiquidPrimary.r, uLiquidPrimary.g, uLiquidPrimary.b, col.a);
    } else if (col.r == 1. && col.g == 0. && col.b == 0.) {
        col = vec4(uLiquidSecondary.r, uLiquidSecondary.g, uLiquidSecondary.b, col.a);
    } else if (col.r == 0. && col.g == 1. && col.b == 0.) {
        col = vec4(uGoop.r, uGoop.g, uGoop.b, col.a);
    }
    fragColor = col;
}
