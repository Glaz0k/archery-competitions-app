import 'dart:collection';
import 'package:flutter/foundation.dart';

class RangeModel with ChangeNotifier {
  final List<Range> _ranges;
  final int maxRanges;
  final int maxShots;
  final bool isShort;

  UnmodifiableListView<Range> get ranges => UnmodifiableListView(_ranges);

  RangeModel(this._ranges, this.maxRanges, this.maxShots, this.isShort);

  void pushRange() {
    _ranges.add(Range());
    notifyListeners();
  }

  void popRange() {
    _ranges.removeLast();
    notifyListeners();
  }
}

class Range with ChangeNotifier {
  final List<ValueNotifier<Shot>> _shots;
  final ValueNotifier<int> _score;

  UnmodifiableListView<ValueListenable<Shot>> get shots =>
      UnmodifiableListView(_shots);

  ValueListenable<int> get score => _score;

  Range._(this._shots, this._score);

  Range() : this._([], ValueNotifier(0));

  void pushShot(Shot shot) {
    _shots.add(ValueNotifier(shot));
    notifyListeners();
  }

  void popShot() {
    _shots.last.dispose();
    _shots.removeLast();
    notifyListeners();
  }

  void changeShot(int index, Shot newValue) {
    Shot old = _shots[index].value;
    _score.value += newValue.algebraicValue - old.algebraicValue;
    _shots[index].value = newValue;
  }
}

class Shot {
  final int rawValue;

  const Shot(this.rawValue) : assert(0 <= rawValue || rawValue <= 11);

  Shot.fromSlider(double sliderValue, bool isShort)
    : this((isShort && sliderValue < 6.0) ? 0 : sliderValue.round());

  @override
  String toString() => switch (rawValue) {
    0 => 'M',
    11 => 'X',
    _ => rawValue.toString(),
  };
  int get algebraicValue => rawValue == 11? 10 : rawValue;
}
