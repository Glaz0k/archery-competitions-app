package com.DevBow.auth.dto;

import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class JwtResponseDto {

    private final String type = "Bearer";
    private String token;

}
