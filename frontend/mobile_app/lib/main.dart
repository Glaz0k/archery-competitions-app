import 'package:flutter/material.dart';
import 'package:mobile_app/api/fake.dart';
import 'package:mobile_app/model/user_model.dart';
import 'package:provider/provider.dart';

import 'api/api.dart';
import 'app.dart';

void main() => runApp(
  MultiProvider(
    providers: [
      ChangeNotifierProvider(create: (context) => UserModel()),
      Provider<Api>(create: (context) => FakeServer()),
    ],
    child: MaterialApp(
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.green),
      ),
      home: OnionApp(), //MainCompetitionPage(),
    ),
  ),
);
