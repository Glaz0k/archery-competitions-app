package com.DevBow.auth.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;

public record GenerateAdminDto(

    @NotBlank(message = "Superuser password must not be blank")
    @JsonProperty("superuser_password")
    String superuserPassword

) {

}
