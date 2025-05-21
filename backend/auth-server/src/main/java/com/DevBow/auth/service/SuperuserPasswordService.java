package com.DevBow.auth.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class SuperuserPasswordService {

    private final PasswordEncoder passwordEncoder;

    private final String encodedSuperuserPassword;

    public SuperuserPasswordService(@Value("${superuser-password}") String superuserPassword, PasswordEncoder encoder) {
        this.passwordEncoder = encoder;
        this.encodedSuperuserPassword = encoder.encode(superuserPassword);
    }

    public boolean isValidSuperuserPassword(String password) {
        return passwordEncoder.matches(password, encodedSuperuserPassword);
    }

}
