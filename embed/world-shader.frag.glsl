#version 330 core

in vec2 vTexCoords;
out vec4 fragColor;
in vec4 vColor;

uniform sampler2D uTexture;
uniform vec4 uTexBounds;

// custom uniforms
uniform int uMode;
uniform float uSpeed;
uniform float uTime;
uniform float uXVar;
uniform float uYVar;
uniform float uCustom;

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
    fragColor = col;
}