package com.DevBow.auth.repository;

import com.DevBow.auth.entity.UserEntity;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;

public interface UserRepository extends JpaRepository< UserEntity, Long > {

    Optional< UserEntity > getUserEntityByLogin(String login);
}
