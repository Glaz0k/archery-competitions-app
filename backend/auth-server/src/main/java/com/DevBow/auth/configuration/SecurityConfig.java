package com.DevBow.auth.configuration;

import com.nimbusds.jose.jwk.JWK;
import com.nimbusds.jose.jwk.OctetSequenceKey;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.convert.converter.Converter;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.jwt.*;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationConverter;
import org.springframework.security.web.SecurityFilterChain;

import javax.crypto.spec.SecretKeySpec;
import java.util.Base64;
import java.util.List;

import static com.nimbusds.jose.JWSAlgorithm.HS256;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Value("${jwt-secret-key}")
    private String jwtSecretKeyBase64;

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http.csrf(AbstractHttpConfigurer::disable)
            .authorizeHttpRequests(authorizeHttpRequests -> authorizeHttpRequests
                .requestMatchers("/api/auth/sign_up").permitAll()
                .requestMatchers("/api/auth/sign_in").permitAll()
                .requestMatchers("/api/auth/generate/admin").permitAll()
                .requestMatchers("/api/auth/generate/user").hasRole("admin")
                .anyRequest().authenticated())
            .oauth2ResourceServer(oauth2 -> oauth2
                .jwt(jwt -> jwt
                    .decoder(jwtDecoder())
                    .jwtAuthenticationConverter(jwtAuthenticationConverter())));

        return http.build();
    }

    @Bean
    public JwtDecoder jwtDecoder() {
        return NimbusJwtDecoder.withSecretKey(new SecretKeySpec(
                getSecretKeyBytes(),
                "HmacSHA256"))
            .build();
    }

    @Bean
    public JwtEncoder jwtEncoder() {
        JWK jwk = new OctetSequenceKey.Builder(getSecretKeyBytes())
            .algorithm(HS256)
            .build();

        return new NimbusJwtEncoder((jwkSelector, securityContext) -> List.of(jwk));
    }

    @Bean
    public Converter< Jwt, AbstractAuthenticationToken > jwtAuthenticationConverter() {
        JwtAuthenticationConverter converter = new JwtAuthenticationConverter();

        converter.setJwtGrantedAuthoritiesConverter(jwt -> {
            String role = jwt.getClaim("role");
            return List.of(new SimpleGrantedAuthority("ROLE_" + role));
        });

        return converter;
    }

    private byte[] getSecretKeyBytes() {
        byte[] secretKeyBytes = Base64.getDecoder().decode(jwtSecretKeyBase64);
        if (secretKeyBytes.length * 8 < 256) {
            throw new IllegalArgumentException("Secret key length must be greater or equal 256 bits");
        }
        return secretKeyBytes;
    }
}
