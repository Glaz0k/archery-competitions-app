package com.DevBow.auth.service;

import com.DevBow.auth.entity.UserEntity;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.io.Encoders;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import javax.crypto.SecretKey;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class JwtProviderTest {

    private final static SecretKey TEST_SECRET_KEY = Jwts.SIG.HS256.key().build();
    @Mock
    private UserEntity user;
    private JwtProvider jwtProvider;

    @BeforeEach
    void setUp() {
        jwtProvider = new JwtProvider(Encoders.BASE64.encode(TEST_SECRET_KEY.getEncoded()));
    }

    @Test
    void testGenerateToken() {
        when(user.getId()).thenReturn(1L);
        when(user.getRole()).thenReturn(UserEntity.Role.ADMIN);

        String token = jwtProvider.generateToken(user);

        assertNotNull(token);

        var claims = Jwts.parser()
            .verifyWith(TEST_SECRET_KEY)
            .build()
            .parseSignedClaims(token)
            .getPayload();

        assertEquals(1L, claims.get("user_id", Long.class));
        assertEquals("admin", claims.get("role", String.class));
    }
}