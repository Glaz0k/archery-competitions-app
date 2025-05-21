package com.DevBow.auth.service;

import com.DevBow.auth.entity.UserEntity;
import lombok.RequiredArgsConstructor;
import org.springframework.lang.NonNull;
import org.springframework.security.oauth2.jose.jws.MacAlgorithm;
import org.springframework.security.oauth2.jwt.JwsHeader;
import org.springframework.security.oauth2.jwt.JwtClaimsSet;
import org.springframework.security.oauth2.jwt.JwtEncoder;
import org.springframework.security.oauth2.jwt.JwtEncoderParameters;
import org.springframework.stereotype.Component;

@RequiredArgsConstructor
@Component
public class JwtProvider {

    private final JwtEncoder jwtEncoder;

    public String generateToken(@NonNull UserEntity user) {
        JwtClaimsSet claims = JwtClaimsSet.builder()
            .claim("user_id", user.getId())
            .claim("role", user.getRole().toString())
            .build();

        JwsHeader header = JwsHeader
            .with(MacAlgorithm.HS256)
            .type("JWT")
            .build();

        return jwtEncoder
            .encode(JwtEncoderParameters.from(header, claims)).getTokenValue();
    }

}
