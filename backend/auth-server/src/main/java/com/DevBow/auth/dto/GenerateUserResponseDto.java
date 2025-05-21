package com.DevBow.auth.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

public record GenerateUserResponseDto(

    @JsonProperty("auth_data")
    AuthDataDto authorizationData,

    @JsonProperty("credentials")
    CredentialsDto credentials

) {

}
