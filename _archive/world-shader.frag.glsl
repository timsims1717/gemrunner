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

// Nearest emulated sample given floating point position and texel offset.
// Also zero's off screen.
vec3 Fetch(vec2 pos,vec2 off) {
    pos=floor(pos*res+off)/res;
    if (max(abs(pos.x-0.5),abs(pos.y-0.5))>0.5) {
        return vec3(0.0,0.0,0.0);
    }
    return texture(uTexture, pos.xy,-16.0).rgb;
}

// Distance in emulated pixels to nearest texel.
vec2 Dist(vec2 pos) {
    pos = pos*res;
    return -((pos-floor(pos))-vec2(0.5));
}

// 1D Gaussian.
float Gaus(float pos, float scale) {
    return exp2(scale*pos*pos);
}

// 3-tap Gaussian filter along horz line.
vec3 Horz3(vec2 pos,float off) {
    vec3 b=Fetch(pos,vec2(-1.0,off));
    vec3 c=Fetch(pos,vec2( 0.0,off));
    vec3 d=Fetch(pos,vec2( 1.0,off));
    float dst=Dist(pos).x;
    // Convert distance to weight.
    float scale=hardPix;
    float wb=Gaus(dst-1.0,scale);
    float wc=Gaus(dst+0.0,scale);
    float wd=Gaus(dst+1.0,scale);
    // Return filtered sample.
    return (b*wb+c*wc+d*wd)/(wb+wc+wd);
}

// 5-tap Gaussian filter along horz line.
vec3 Horz5(vec2 pos,float off) {
    vec3 a=Fetch(pos,vec2(-2.0,off));
    vec3 b=Fetch(pos,vec2(-1.0,off));
    vec3 c=Fetch(pos,vec2( 0.0,off));
    vec3 d=Fetch(pos,vec2( 1.0,off));
    vec3 e=Fetch(pos,vec2( 2.0,off));
    float dst=Dist(pos).x;
    // Convert distance to weight.
    float scale=hardPix;
    float wa=Gaus(dst-2.0,scale);
    float wb=Gaus(dst-1.0,scale);
    float wc=Gaus(dst+0.0,scale);
    float wd=Gaus(dst+1.0,scale);
    float we=Gaus(dst+2.0,scale);
    // Return filtered sample.
    return (a*wa+b*wb+c*wc+d*wd+e*we)/(wa+wb+wc+wd+we);
}

// Return scanline weight.
float Scan(vec2 pos,float off) {
    float dst = Dist(pos).y;
    return Gaus(dst+off,hardScan);
}

// Allow nearest three lines to effect pixel.
vec3 Tri(vec2 pos) {
    vec3 a=Horz3(pos,-1.0);
    vec3 b=Horz5(pos, 0.0);
    vec3 c=Horz3(pos, 1.0);
    float wa=Scan(pos,-1.0);
    float wb=Scan(pos, 0.0);
    float wc=Scan(pos, 1.0);
    return a*wa+b*wb+c*wc;
}

// Distortion of scanlines, and end of screen alpha.
vec2 Warp(vec2 pos) {
    pos = pos*2.0-1.0;
    pos *= vec2(1.0+(pos.y*pos.y)*warp.x, 1.0+(pos.x*pos.x)*warp.y);
    return pos*0.5+0.5;
}

// Shadow mask.
vec3 Mask(vec2 pos) {
    pos.x += pos.y*3.0;
    vec3 mask = vec3(maskDark, maskDark, maskDark);
    pos.x = fract(pos.x/6.0);
    if (pos.x<0.333) {
        mask.r = maskLight;
    } else if (pos.x<0.666) {
        mask.g = maskLight;
    } else {
        mask.b = maskLight;
    }
    return mask;
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