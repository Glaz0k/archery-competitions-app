import 'dart:io';

class NotFoundException extends HttpException {
  NotFoundException(super.message);
}

class InvalidParametersException extends HttpException {
  InvalidParametersException(super.message);
}

class AlreadyExistException extends HttpException {
  AlreadyExistException(super.message);
}

class BadActionException extends HttpException {
  BadActionException(super.message);
}