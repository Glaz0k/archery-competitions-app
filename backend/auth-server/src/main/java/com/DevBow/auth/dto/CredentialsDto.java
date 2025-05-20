package com.DevBow.auth.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;

public record CredentialsDto(

    @NotNull(message = "Login must be non-null")
    @Pattern(regexp = "^[a-zA-Z0-9._-]{6,20}$", message = "Invalid login")
    @JsonProperty("login")
    String login,

    @NotNull(message = "Password must be non-null")
    @Pattern(regexp = "^[a-zA-Z0-9._-]{6,20}$", message = "Invalid password")
    @JsonProperty("password")
    String password

) {

}