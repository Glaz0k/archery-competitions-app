package com.DevBow.auth.controller;

import com.DevBow.auth.dto.CredentialsDto;
import com.DevBow.auth.dto.JwtResponseDto;
import com.DevBow.auth.entity.UserEntity;
import com.DevBow.auth.repository.UserRepository;
import com.DevBow.auth.service.JwtProvider;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

@Slf4j
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/auth")
public class AuthenticationController {

    private final UserRepository userRepository;

    private final PasswordEncoder passwordEncoder;

    private final JwtProvider jwtProvider;

    @Transactional
    @PostMapping("/registration")
    public JwtResponseDto registerUser(@RequestBody @Validated CredentialsDto credentials) {
        log.debug("POST /api/auth/registration method invoked");

        UserEntity user = UserEntity.builder()
            .login(credentials.login())
            .password(passwordEncoder.encode(credentials.password()))
            .role(UserEntity.Role.USER)
            .build();
        log.debug("New user built");

        try {
            user = userRepository.save(user);
            userRepository.flush();
        } catch (DataIntegrityViolationException e) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "EXISTS", e);
        }
        log.debug("New user registered");

        return new JwtResponseDto(jwtProvider.generateToken(user));
    }

    @Transactional
    @PostMapping("/login")
    public JwtResponseDto loginUser(@RequestBody @Validated CredentialsDto credentials) {
        log.debug("POST /api/auth/login method invoked");

        UserEntity user = userRepository.getUserEntityByLogin(credentials.login())
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.BAD_REQUEST, "INVALID PARAMETERS"));
        log.debug("User found");

        if (!passwordEncoder.matches(credentials.password(), user.getPassword())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "INVALID PARAMETERS");
        }
        log.debug("Credentials matches");

        return new JwtResponseDto(jwtProvider.generateToken(user));
    }

}

