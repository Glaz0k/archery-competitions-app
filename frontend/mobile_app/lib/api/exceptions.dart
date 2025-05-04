import 'dart:io';

import 'package:mobile_app/api/common.dart';

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

class InvalidScoreException extends HttpException {
  final int shotOrdinal;
  final RangeType rangeType;
  InvalidScoreException(super.message, this.shotOrdinal, this.rangeType);
}
