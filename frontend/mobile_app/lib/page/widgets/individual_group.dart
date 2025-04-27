import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'user.dart';

class IndividualGroup extends StatefulWidget {

  const IndividualGroup({super.key});

  @override
  State<StatefulWidget> createState() => _IndividualGroup();
}

class _IndividualGroup extends State<IndividualGroup> {
  final List<BowClass> combinedGroups = [BowClass.classic, BowClass.block];

  late User _user;
  @override
  void initState() {
    super.initState();
    _user = Provider.of<UserProvider>(context, listen: false).userPref;
  }

  @override
  Widget build(BuildContext context) {
    final user = Provider.of<UserProvider>(context).userPref;
    return Scaffold(
      appBar: AppBar(
        leading: BackButton(),
        title: Text("Индивидуальные группы", style: Theme.of(context).appBarTheme.titleTextStyle,),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            buildGroup(),
            Padding(padding: EdgeInsets.all(8)),
            if(_user.bow != null && combinedGroups.contains(_user.bow!))
              buildCombinedGroup(),
          ],
        ),
      ),
    );
  }

  Widget buildGroup() {
    final gender = _user.identity;
    final bow = _user.bow?.getBowClass;
    final type = gender == Gender.male ? "Мужчины" : "Женщины";
    return ElevatedButton(onPressed: () {},
        child: Text('$bow - $type', style: TextStyle(fontSize: 19, fontWeight: FontWeight.w900),));
  }

  Widget buildCombinedGroup() {
    final bow = _user.bow?.getBowClass;
    return ElevatedButton(onPressed: () {},
        child: Text("$bow - Объединенные", style: TextStyle(fontSize: 19, fontWeight: FontWeight.w900),)
    );
  }
}