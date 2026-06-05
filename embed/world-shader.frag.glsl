#version 330 core

in vec2 vTexCoords;
out vec4 fragColor;
in vec4 vColor;

uniform sampler2D uTexture;
uniform vec4 uTexBounds;

// custom uniforms
uniform int uDarkness;
uniform float uDarknessDist;
uniform float uDarknessWidth;
uniform float uDarknessGrad;
uniform vec2 uPlayer1Loc;
uniform vec2 uPlayer2Loc;
uniform vec2 uPlayer3Loc;
uniform vec2 uPlayer4Loc;
uniform int uMode;
uniform float uSpeed;
uniform float uTime;
uniform float uXVar;
uniform float uYVar;
uniform float uCustom;

// color uniforms
uniform vec3 uPrimary;
uniform vec3 uSecondary;
uniform vec3 uDoodad;
uniform vec3 uGoop;
uniform vec3 uLiquidPrimary;
uniform vec3 uLiquidSecondary;

// Hardness of scanline.
//  -8.0 = soft
// -16.0 = medium
float hardScan=-8.0;

// Hardness of pixels in scanline.
// -2.0 = soft
// -4.0 = hard
float hardPix=-3.0;

// Display warp.
// 0.0 = none
// 1.0/8.0 = extreme
vec2 warp=vec2(1.0/32.0,1.0/24.0);

// Amount of shadow mask.
float maskDark=0.5;
float maskLight=1.5;

#define res (uTexBounds.zw)

float rand(vec2 n) {
    return fract(sin(dot(n, vec2(12.9898, 4.1414))) * 43758.5453);
}

float noise(vec2 n) {
    const vec2 d = vec2(0.0, 1.0);
    vec2 b = floor(n), f = smoothstep(vec2(0.0), vec2(1.0), fract(n));
    return mix(mix(rand(b), rand(b + d.yx), f.x), mix(rand(b + d.xy), rand(b + d.yy), f.x), f.y);
}

void main() {
    vec2 t = vTexCoords / uTexBounds.zw;
    vec4 col;
    //    vec2 pos = Warp(vTexCoords / uTexBounds.zw);
    //    col.rgb = Tri(pos)*Mask(vTexCoords.xy);
    switch (uMode) {
    case 0: // none
        col = texture(uTexture, t).rgba;
        break;
    case 1: // watery
        t.y += cos(t.x * uCustom + (uTime * uSpeed))*uXVar;
        t.x += cos(t.y * uCustom + (uTime * uSpeed))*uYVar;
        col = texture(uTexture, t).rgba;
        break;
    case 2: // heat
        vec2 p_d = t;
        p_d.y -= uTime * uSpeed;
        vec4 dst_map_val = vec4(noise(p_d * vec2(50)));
        vec2 dst_offset = dst_map_val.xy;
        dst_offset -= vec2(.5,.5);
        dst_offset *= uCustom;
        dst_offset *= 0.01;
        dst_offset *= min(uYVar - t.t, 1.);
        vec2 dist_tex_coord = t.st + dst_offset;
        col = texture(uTexture, dist_tex_coord).rgba;
        break;
    }

    // colors
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

    // darkness
    if (uDarkness == 1) {
        float ar = (uTexBounds.z / uTexBounds.w) * uDarknessWidth;
        vec2 ot = vTexCoords / uTexBounds.zw;
        ot.x -= 0.5;
        ot.x *= ar;
        vec4 shadowCol = vec4(0., 0., 0., 1.);
        float g = 1;
        for (int i = 0; i < 4; i++) {
            vec2 pLoc;
            switch (i) {
            case 0:
                pLoc = uPlayer1Loc;
                break;
            case 1:
                pLoc = uPlayer2Loc;
                break;
            case 2:
                pLoc = uPlayer3Loc;
                break;
            case 3:
                pLoc = uPlayer4Loc;
                break;
            }
            if (pLoc.x < 0 || pLoc.y < 0 || pLoc.x > 1 || pLoc.y > 1) {
                continue;
            }
            pLoc.x -= 0.5;
            pLoc.x *= ar;
            float dist = abs(distance(pLoc, ot));
            float dg = 1;
            if (dist < uDarknessDist) {
                dg = 0;
            } else if (dist > uDarknessDist + uDarknessGrad) {
                dg = 1;
            } else {
                float d = uDarknessDist + uDarknessGrad - dist;
                dg = 1-d/uDarknessGrad;
            }
            if (dg < g) {
                g = dg;
            }
            g = clamp(g, 0, 1);
        }
        col = mix(col, shadowCol, g);
    }
    fragColor = col;
}