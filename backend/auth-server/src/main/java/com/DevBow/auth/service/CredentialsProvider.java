package com.DevBow.auth.service;

import com.DevBow.auth.dto.CredentialsDto;
import org.apache.commons.text.RandomStringGenerator;
import org.springframework.stereotype.Component;

@Component
public class CredentialsProvider {

    private static final int MIN_LENGTH = 6;
    private static final int MAX_LENGTH = 20;
    private static final char[][] RANGES = {
        { 'a', 'z' },
        { 'A', 'Z' },
        { '0', '9' },
        { '.', '.' },
        { '_', '_' },
        { '-', '-' }
    };

    private static final RandomStringGenerator generator = new RandomStringGenerator.Builder()
        .withinRange(RANGES)
        .get();

    private String generateCredential() {
        return generator.generate(MIN_LENGTH, MAX_LENGTH);
    }

    public CredentialsDto generateCredentials() {
        return new CredentialsDto(
            generateCredential(),
            generateCredential()
        );
    }
}
