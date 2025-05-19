import 'package:flutter/material.dart';
import 'package:mobile_app/api/fake.dart';
import 'package:mobile_app/page/edit_profile_page.dart';
import 'package:mobile_app/page/main_competition_page.dart';
import 'package:mobile_app/page/widgets/user.dart';
import 'package:provider/provider.dart';

import 'api/api.dart';
import 'model/range_model.dart';
import 'app.dart';
void main() => runApp(
  MultiProvider(
    providers: [
      ChangeNotifierProvider(create: (context) => UserProvider()),
      Provider<Api>(create: (context) => FakeServer()),
    ],
    child: MaterialApp(
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.pinkAccent),
      ),
      home: OnionApp(), //MainCompetitionPage(),
    ),
  ),
);
