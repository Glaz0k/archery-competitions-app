package com.DevBow.auth.dto;

import com.DevBow.auth.entity.UserEntity;
import com.fasterxml.jackson.annotation.JsonProperty;

public record AuthDataDto(

    @JsonProperty("user_id")
    Long userId,

    @JsonProperty("role")
    UserEntity.Role userRole

) {

}
