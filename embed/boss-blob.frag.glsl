#version 330 core

in vec2  vTexCoords;
out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

uniform float uTime;
uniform int uBallCount;
uniform vec2 uHeadPos;

uniform vec3 uColorInner;
uniform vec3 uColorOuter;

float mBall(vec2 uv, vec2 pos, vec2 scale, float radius) {
    uv -= pos;
    uv *= scale;
    return radius/(dot(uv,uv));
}

// from https://www.shadertoy.com/view/MllXDH
void main() {
    // Get our current screen coordinate
    vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
    vec4 col = texture(uTexture, t);

    vec3 color_inner = uColorInner;
    vec3 color_outer = uColorOuter;

    vec2 s = uTexBounds.xy;
    vec2 uv = (vTexCoords-s)/uTexBounds.zw;

    float xd = uHeadPos.x - 0.5;
    float xi = xd / uBallCount;
    float yd = uHeadPos.y;
    float yi = yd / uBallCount;

    float mb = 0.0;
    vec2 headScale = vec2(1.25, 1.);
    vec2 neckScale = vec2(2., 1.);
    mb += mBall(uv, uHeadPos, headScale, 0.01); // head
    for (int i = 1; i <= uBallCount; i++) { // neck
        mb += mBall(uv, uHeadPos - vec2(xi, yi)*i, neckScale, 0.005);
    }
    mb += mBall(uv, vec2(0.5, -0.09), vec2(0.85, 1.), 0.02); // base

    vec3 mbext = color_outer * (1.-smoothstep(mb, mb+0.01, 0.5)); // 0.5 for control the blob thickness
    vec3 mbin = color_inner * (1.-smoothstep(mb, mb+0.01, 0.8));  // 0.8 for control the blob kernel size

    fragColor.rgb = max(mbin, mbext);
    if (fragColor.r == 0. && fragColor.g == 0. && fragColor.b == 0.) {
        fragColor = col;
    }
}
