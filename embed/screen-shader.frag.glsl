#version 330 core

in vec2 vTexCoords;
out vec4 fragColor;
in vec4 vColor;

uniform sampler2D uTexture;
uniform vec4 uTexBounds;

// custom uniforms
uniform float uSpeed;
uniform float uTime;

#define res (uTexBounds.zw)

// Hardness of scanline.
//  -8.0 = soft
// -16.0 = medium
float hardScan = -8.0;

// Hardness of pixels in scanline.
// -2.0 = soft
// -4.0 = hard
float hardPix = -3.0;

// Display warp.
// 0.0 = none
// 1.0/8.0 = extreme
vec2 warp = vec2(1.0/32.0, 1.0/24.0);

// Amount of shadow mask.
float maskDark = 0.5;
float maskLight = 1.5;

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
    vec3 b = Fetch(pos,vec2(-1.0, off));
    vec3 c = Fetch(pos,vec2(0.0, off));
    vec3 d = Fetch(pos,vec2(1.0, off));
    float dst = Dist(pos).x;
    // Convert distance to weight.
    float scale = hardPix;
    float wb = Gaus(dst-1.0, scale);
    float wc = Gaus(dst+0.0, scale);
    float wd = Gaus(dst+1.0, scale);
    // Return filtered sample.
    return (b*wb+c*wc+d*wd)/(wb+wc+wd);
}

// 5-tap Gaussian filter along horz line.
vec3 Horz5(vec2 pos,float off) {
    vec3 a = Fetch(pos, vec2(-2.0, off));
    vec3 b = Fetch(pos, vec2(-1.0, off));
    vec3 c = Fetch(pos, vec2(0.0, off));
    vec3 d = Fetch(pos, vec2(1.0, off));
    vec3 e = Fetch(pos, vec2(2.0, off));
    float dst=Dist(pos).x;
    // Convert distance to weight.
    float scale = hardPix;
    float wa = Gaus(dst-2.0, scale);
    float wb = Gaus(dst-1.0, scale);
    float wc = Gaus(dst+0.0, scale);
    float wd = Gaus(dst+1.0, scale);
    float we = Gaus(dst+2.0, scale);
    // Return filtered sample.
    return (a*wa+b*wb+c*wc+d*wd+e*we)/(wa+wb+wc+wd+we);
}

// Return scanline weight.
float Scan(vec2 pos,float off) {
    float dst = Dist(pos).y;
    return Gaus(dst+off, hardScan);
}

// Allow nearest three lines to effect pixel.
vec3 Tri(vec2 pos) {
    vec3 a = Horz3(pos, -1.0);
    vec3 b = Horz5(pos, 0.0);
    vec3 c = Horz3(pos, 1.0);
    float wa = Scan(pos, -1.0);
    float wb = Scan(pos, 0.0);
    float wc = Scan(pos, 1.0);
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
    vec3 col;
    vec2 pos = Warp(vTexCoords / uTexBounds.zw);
    hardScan = -12.0;
    maskDark = maskLight = 1.0;
    col.rgb = Tri(pos)*Mask(vTexCoords.xy);
    fragColor = vec4(col.rgb, 1.0);
}