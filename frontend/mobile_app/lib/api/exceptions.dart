abstract class OnionException implements Exception {
  final String message;

  OnionException(this.message);
}

class NotFoundException extends OnionException {
  NotFoundException(super.message);
}

class InvalidParametersException extends OnionException {
  InvalidParametersException(super.message);
}

class AlreadyExistException extends OnionException {
  AlreadyExistException(super.message);
}

class BadActionException extends OnionException {
  BadActionException(super.message);
}

class InvalidScoreException extends OnionException {
  final int shotOrdinal;
  final String rangeType;

  InvalidScoreException(super.message, this.shotOrdinal, this.rangeType);
}
