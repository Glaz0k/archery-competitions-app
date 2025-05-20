import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:mobile_app/model/user_model.dart';
import 'package:mobile_app/page/login_page.dart';
import 'package:mobile_app/page/main_competition_page.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

import 'api/api.dart';

class OnionApp extends StatelessWidget {
  const OnionApp({super.key});

  @override
  Widget build(BuildContext context) {
    var userProvider = context.watch<UserModel>();
    var api = context.watch<Api>();
    loadingScreen() => Scaffold(
      appBar: OnionBar.withoutProfile("Загрузка...", context),
      body: Center(child: CircularProgressIndicator()),
    );
    if (userProvider.loading) {
      log("Ещё не загрузилось хранилище");
      return loadingScreen();
    } else {
      if (userProvider.getId() == null) {
        log("Нету id");
        return LoginPage();
      } else {
        if (userProvider.user == null) {
          userProvider.loadUser(api);
          return loadingScreen();
        } else {
          return MainCompetitionPage();
        }
      }
    }
  }
}
