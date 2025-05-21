package com.DevBow.auth.controller;

import com.DevBow.auth.dto.*;
import com.DevBow.auth.entity.UserEntity;
import com.DevBow.auth.repository.UserRepository;
import com.DevBow.auth.service.CredentialsProvider;
import com.DevBow.auth.service.JwtProvider;
import com.DevBow.auth.service.SuperuserPasswordService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
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
@Transactional
@RequestMapping("/api/auth")
public class AuthenticationController {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtProvider jwtProvider;
    private final SuperuserPasswordService superuserPasswordService;
    private final CredentialsProvider credentialsProvider;

    @PostMapping("/sign_up")
    public JwtAuthResponseDto signUpUser(@RequestBody @Validated CredentialsDto credentials) {
        log.info("Sign up attempt, login: {}", credentials.login());

        UserEntity user = UserEntity.builder()
            .login(credentials.login())
            .password(passwordEncoder.encode(credentials.password()))
            .role(UserEntity.Role.USER)
            .build();

        try {
            user = userRepository.saveAndFlush(user);
        } catch (DataIntegrityViolationException e) {
            throw new ResponseStatusException(
                HttpStatus.BAD_REQUEST,
                "USER WITH LOGIN EXISTS: %s".formatted(credentials.login()),
                e);
        }

        log.info("Successfully signed up, id: {}", user.getId());
        return new JwtAuthResponseDto(
            new AuthDataDto(
                user.getId(),
                user.getRole()
            ),
            jwtProvider.generateToken(user)
        );
    }

    @PostMapping("/sign_in")
    public JwtAuthResponseDto signInUser(@RequestBody @Validated CredentialsDto credentials) {
        log.info("Sign in attempt, login: {}", credentials.login());

        UserEntity user = userRepository
            .getUserEntityByLogin(credentials.login())
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.BAD_REQUEST, "INVALID CREDENTIALS"));

        if (!passwordEncoder.matches(credentials.password(), user.getPassword())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "INVALID CREDENTIALS");
        }

        log.debug("Successfully signed in, id: {}", user.getId());
        return new JwtAuthResponseDto(
            new AuthDataDto(
                user.getId(),
                user.getRole()
            ),
            jwtProvider.generateToken(user)
        );
    }

    @PostMapping("/generate/admin")
    public GenerateUserResponseDto generateAdmin(@RequestBody @Validated GenerateAdminDto generateAdminDto) {
        log.info("Generate admin attempt");

        if (!superuserPasswordService.isValidSuperuserPassword(generateAdminDto.superuserPassword())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "INVALID SUPERUSER PASSWORD");
        }

        CredentialsDto credentials = credentialsProvider.generateCredentials();

        UserEntity generatedAdmin = UserEntity.builder()
            .login(credentials.login())
            .password(passwordEncoder.encode(credentials.password()))
            .role(UserEntity.Role.ADMIN)
            .build();

        generatedAdmin = userRepository.saveAndFlush(generatedAdmin);

        log.info("Successfully generated admin account, id: {}", generatedAdmin.getId());
        return new GenerateUserResponseDto(
            new AuthDataDto(
                generatedAdmin.getId(),
                generatedAdmin.getRole()
            ),
            credentials
        );
    }

    @PostMapping("/generate/user")
    public GenerateUserResponseDto generateUser(JwtAuthenticationToken authentication) {
        log.info("Generate user attempt, initiated by admin, id: {}",
            (long) authentication.getToken().getClaim("user_id"));

        CredentialsDto credentials = credentialsProvider.generateCredentials();

        UserEntity generatedUser = UserEntity.builder()
            .login(credentials.login())
            .password(passwordEncoder.encode(credentials.password()))
            .role(UserEntity.Role.USER)
            .build();

        generatedUser = userRepository.saveAndFlush(generatedUser);

        log.info("Successfully generated user account, id: {}", generatedUser.getId());
        return new GenerateUserResponseDto(
            new AuthDataDto(
                generatedUser.getId(),
                generatedUser.getRole()
            ),
            credentials
        );
    }
}

