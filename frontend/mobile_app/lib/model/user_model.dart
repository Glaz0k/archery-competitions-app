import 'package:flutter/material.dart';
import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/api/responses.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../api/api.dart';

const String userKey = "user";

class UserModel with ChangeNotifier {
  bool loading = true;
  SharedPreferencesWithCache? _prefs;
  CompetitorFull? user;

  UserModel() {
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

  Future<CompetitorFull> loadUser(Api api) async {
    user = await api.getCompetitor(getId()!);
    notifyListeners();
    return user!;
  }

  Future<CompetitorFull> setUser(Api api, ChangeCompetitor request) async {
    user = await api.putCompetitor(getId()!, request);
    notifyListeners();
    return user!;
  }
}
