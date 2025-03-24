package com.DevBow.auth.service;

import com.DevBow.auth.entity.UserEntity;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.io.Decoders;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;

@Component
public class JwtProvider {

    private final SecretKey jwtSecretKey;

    public JwtProvider(@Value("${jwt.secret-key}") String jwtSecretKeyBase64) {
        this.jwtSecretKey = Keys.hmacShaKeyFor(Decoders.BASE64.decode(jwtSecretKeyBase64));
    }

    public String generateToken(@NonNull UserEntity user) {
        return Jwts.builder()
            .claim("user_id", user.getId())
            .claim("role", user.getRole().toString())
            .signWith(jwtSecretKey, Jwts.SIG.HS256)
            .compact();
    }

}
