package com.DevBow.auth.controller;

import com.DevBow.auth.dto.CredentialsDto;
import com.DevBow.auth.dto.JwtResponseDto;
import com.DevBow.auth.entity.UserEntity;
import com.DevBow.auth.repository.UserRepository;
import com.DevBow.auth.service.JwtProvider;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.server.ResponseStatusException;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AuthenticationControllerTest {

    @Mock
    private UserRepository userRepository;

    @Mock
    private PasswordEncoder passwordEncoder;

    @Mock
    private JwtProvider jwtProvider;

    @InjectMocks
    private AuthenticationController authenticationController;

    private CredentialsDto validCredentials;
    private CredentialsDto invalidCredentials;
    private UserEntity userEntity;
    private String encodedPassword;
    private String jwtToken;

    @BeforeEach
    void setUp() {
        validCredentials = new CredentialsDto("testUser", "testPass");
        invalidCredentials = new CredentialsDto("wrongUser", "wrongPass");
        encodedPassword = "encodedTestPass";
        jwtToken = "test.jwt.token";

        userEntity = UserEntity.builder()
            .login(validCredentials.login())
            .password(encodedPassword)
            .role(UserEntity.Role.USER)
            .build();
    }

    @Test
    void registerUserSuccess() {
        when(passwordEncoder.encode(anyString())).thenReturn(encodedPassword);
        when(userRepository.save(any(UserEntity.class))).thenReturn(userEntity);
        when(jwtProvider.generateToken(any(UserEntity.class))).thenReturn(jwtToken);

        JwtResponseDto response = authenticationController.registerUser(validCredentials);

        assertNotNull(response);
        assertEquals("Bearer", response.getType());
        assertEquals(jwtToken, response.getToken());

        verify(userRepository, times(1)).save(any(UserEntity.class));
        verify(jwtProvider, times(1)).generateToken(any(UserEntity.class));
    }

    @Test
    void registerUserAlreadyExists() {
        when(passwordEncoder.encode(anyString())).thenReturn(encodedPassword);
        when(userRepository.save(any(UserEntity.class)))
            .thenThrow(new DataIntegrityViolationException("User already exists"));

        ResponseStatusException exception = assertThrows(ResponseStatusException.class,
            () -> authenticationController.registerUser(validCredentials));

        assertEquals(HttpStatus.BAD_REQUEST, exception.getStatusCode());
        assertEquals("EXISTS", exception.getReason());

        verify(userRepository, times(1)).save(any(UserEntity.class));
    }

    @Test
    void loginUserSuccess() {
        when(userRepository.getUserEntityByLogin(validCredentials.login()))
            .thenReturn(Optional.of(userEntity));
        when(passwordEncoder.matches(validCredentials.password(), encodedPassword))
            .thenReturn(true);
        when(jwtProvider.generateToken(userEntity)).thenReturn(jwtToken);

        JwtResponseDto response = authenticationController.loginUser(validCredentials);

        assertNotNull(response);
        assertEquals("Bearer", response.getType());
        assertEquals(jwtToken, response.getToken());

        verify(userRepository, times(1)).getUserEntityByLogin(validCredentials.login());
        verify(passwordEncoder, times(1)).matches(validCredentials.password(), encodedPassword);
    }

    @Test
    void loginUserNotFound() {
        when(userRepository.getUserEntityByLogin(invalidCredentials.login()))
            .thenReturn(Optional.empty());

        ResponseStatusException exception = assertThrows(ResponseStatusException.class,
            () -> authenticationController.loginUser(invalidCredentials));

        assertEquals(HttpStatus.BAD_REQUEST, exception.getStatusCode());
        assertEquals("INVALID PARAMETERS", exception.getReason());

        verify(userRepository, times(1)).getUserEntityByLogin(invalidCredentials.login());
    }

    @Test
    void loginUserInvalidPassword() {
        when(userRepository.getUserEntityByLogin(validCredentials.login()))
            .thenReturn(Optional.of(userEntity));
        when(passwordEncoder.matches(validCredentials.password(), encodedPassword))
            .thenReturn(false);

        ResponseStatusException exception = assertThrows(ResponseStatusException.class,
            () -> authenticationController.loginUser(validCredentials));

        assertEquals(HttpStatus.BAD_REQUEST, exception.getStatusCode());
        assertEquals("INVALID PARAMETERS", exception.getReason());

        verify(userRepository, times(1)).getUserEntityByLogin(validCredentials.login());
        verify(passwordEncoder, times(1)).matches(validCredentials.password(), encodedPassword);
    }
}