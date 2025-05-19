import 'package:flutter/material.dart';
import 'package:mobile_app/page/final_page.dart';
import 'package:mobile_app/page/group_competitors_page.dart';
import 'package:mobile_app/page/qualification_page.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';

class HubPage extends StatelessWidget {
  final int individualGroupId;
  final String title;

  const HubPage({
    super.key,
    required this.individualGroupId,
    required this.title,
  });

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      initialIndex: 1,
      length: 3,
      child: Scaffold(
        appBar: OnionBar(
          title,
          context,
          bottom: TabBar(
            tabs: [
              Tab(text: "Участники"),
              Tab(text: "Секция"),
              Tab(text: "Финал"),
            ],
          ),
        ),
        body: TabBarView(
          children: [
            GroupCompetitorsPage(groupId: individualGroupId),
            QualificationPage(sectionId: 0),
            // TODO: У этой страницы не должно быть параметра
            FinalPage(groupId: individualGroupId),
          ],
        ),
      ),
    );
  }
}
