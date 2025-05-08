import 'package:flutter/material.dart';
import 'package:mobile_app/api/fake.dart';
import 'package:mobile_app/page/main_competition_page.dart';
import 'package:mobile_app/page/widgets/user.dart';
import 'package:provider/provider.dart';

import 'model/range_model.dart';

void main() => runApp(

  MultiProvider(
    providers: [
      ChangeNotifierProvider(
        create: (context) => RangeModel([], 3, 3, false),
      ),
      ChangeNotifierProvider(create: (context) => UserProvider()),
      Provider(create: (context) => FakeServer()),
    ],
    child: MaterialApp(
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.pinkAccent),
      ),
      home: MainCompetitionPage(),
    ),
  ),
);
