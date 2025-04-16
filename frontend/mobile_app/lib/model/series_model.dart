import 'dart:collection';
import 'dart:developer';
import 'package:flutter/foundation.dart';

class SeriesModel extends ChangeNotifier {
  late final List<Series> _seriesList;
  UnmodifiableListView<Series> get seriesList =>
      UnmodifiableListView(_seriesList);
  SeriesModel() {
    _seriesList = [];
  }
  void addToLast(int score) {
    _seriesList.last._scores.add(score);
    notifyListeners();
  }
  void changeElementInLast(int index, int score) {
    _seriesList.last._scores[index] = score;
    notifyListeners();
  }
  void addNewSeries() {
    log('Added new series. All series: $_seriesList');
    _seriesList.add(Series([]));
    notifyListeners();
  }
}

class Series {
  final List<int> _scores;
  UnmodifiableListView<int> get scores => UnmodifiableListView(_scores);
  Series(this._scores);
}