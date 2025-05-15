import 'package:flutter/foundation.dart';
import 'package:mobile_app/api/exceptions.dart';
import 'package:mobile_app/api/responses.dart';

import '../api/requests.dart';

class RangeModel with ChangeNotifier {
  final Future<RangeGroup> Function() _getRangeGroup;
  final Future<Range> Function(ChangeRange) _putRange;
  final Future<Range> Function(int) _endRange;

  RangeGroup rangeGroup;

  RangeModel._(this.rangeGroup, this._getRangeGroup, this._putRange, this._endRange);
  RangeModel(
    RangeGroup rangeGroup, {
    required Future<RangeGroup> Function() getRangeGroup,
    required Future<Range> Function(ChangeRange) putRange,
    required Future<Range> Function(int) endRange,
  }) : this._(rangeGroup, getRangeGroup, putRange, endRange);

  Future<void> registerAndEndRange(int rangeIdx) async {
    var range = rangeGroup.ranges[rangeIdx-1];
    if (range.shots!.any((shot) => shot.score == null)) {
      throw BadActionException("Нельзя закончить незаконченную серию");
    }
    await _putRange(ChangeRange(rangeIdx,range.shots));
    await _endRange(rangeIdx);
    // Ибо сказал Господь: рефетчи.
    rangeGroup = await _getRangeGroup();
    notifyListeners();
  }
  Future<void> reloadRangeGroup() async {
    rangeGroup = await _getRangeGroup();
    notifyListeners();
  }
}