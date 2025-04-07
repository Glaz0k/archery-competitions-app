package com.DevBow.auth.exceptions;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.server.ResponseStatusException;

@ControllerAdvice
public class RestExceptionHandler {

    @ExceptionHandler(ResponseStatusException.class)
    public ResponseEntity< ErrorResponse > handleResponseStatusException(ResponseStatusException ex) {
        return new ResponseEntity<>(new ErrorResponse(ex.getReason()), ex.getStatusCode());
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity< ErrorResponse > handleMethodArgumentNotValidException(MethodArgumentNotValidException ex) {
        return new ResponseEntity<>(new ErrorResponse("INVALID PARAMETERS"), ex.getStatusCode());
    }

    @ExceptionHandler(HttpMessageNotReadableException.class)
    public ResponseEntity< ErrorResponse > handleHttpMessageNotReadableException() {
        return new ResponseEntity<>(new ErrorResponse("INVALID PARAMETERS"), HttpStatus.BAD_REQUEST);
    }
}