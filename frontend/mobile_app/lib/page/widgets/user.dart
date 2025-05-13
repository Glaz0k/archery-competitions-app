import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/api/responses.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../api/api.dart';

const String userKey = "user";
class UserProvider with ChangeNotifier {
  bool loading = true;
  SharedPreferencesWithCache? _prefs;
  CompetitorFull? _user;
  UserProvider() {
    SharedPreferencesWithCache.create(
      cacheOptions: const SharedPreferencesWithCacheOptions(
        allowList: {userKey},
      ),
    ).then((prefs) {
      _prefs = prefs;
      loading = false;
      notifyListeners();
    });
  }
  int? getId() {
    return _prefs?.getInt(userKey);
  }
  void setId(int userId) {
    _prefs?.setInt(userKey, userId).then((_) => notifyListeners());
  }
  void clearId() {
    _prefs?.remove(userKey).then((_) => notifyListeners());
  }

  CompetitorFull? getUser(Api api) {
    return _user;
  }

  Future<CompetitorFull> loadUser(Api api) async {
    var user = await api.getCompetitor(getId()!);
    _user = user;
    notifyListeners();
    return user;
  }

  Future<CompetitorFull> setUser(Api api, ChangeCompetitor request) async {
    var user = await api.putCompetitor(getId()!, request);
    _user = user;
    notifyListeners();
    return user;
  }
}