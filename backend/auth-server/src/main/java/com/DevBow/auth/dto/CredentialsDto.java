package com.DevBow.auth.dto;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

public record CredentialsDto(

    @NotNull
    @Size(
        message = "Invalid login size",
        min = 6,
        max = 20)
    @Pattern(
        message = "Invalid login characters",
        regexp = "^[a-zA-Z0-9._-]+$")
    String login,

    @NotNull
    @Size(
        message = "Invalid password size",
        min = 6,
        max = 20)
    @Pattern(
        message = "Invalid password characters",
        regexp = "^[a-zA-Z0-9._-]+$")
    String password) {

}