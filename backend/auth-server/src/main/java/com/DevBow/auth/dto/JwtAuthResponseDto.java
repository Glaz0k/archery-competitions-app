package com.DevBow.auth.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

public record JwtAuthResponseDto(

    @JsonProperty("auth_data")
    AuthDataDto authData,

    @JsonProperty("token")
    String jwtToken

) {

}
