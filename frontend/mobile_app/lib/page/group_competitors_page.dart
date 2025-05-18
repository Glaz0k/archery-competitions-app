import 'package:flutter/material.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

import '../api/api.dart';
import '../api/responses.dart';

class GroupCompetitorsPage extends StatefulWidget {
  final int groupId;
  const GroupCompetitorsPage({super.key, required this.groupId});

  @override
  State<GroupCompetitorsPage> createState() => _GroupCompetitorsPageState();
}

class _GroupCompetitorsPageState extends State<GroupCompetitorsPage> {
  List<CompetitorGroupDetail> _competitorsFuture = [];

  @override
  void initState() {
    super.initState();
    _loadCompetitors();
  }

  Future<void> _loadCompetitors() async {
    var api = context.read<Api>();
    final competitors = await api.getIndividualGroupCompetitors(widget.groupId);
    if (competitors.isEmpty) return;
    setState(() {
      _competitorsFuture = competitors;
    });
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: _loadCompetitors,
      child: Scaffold(
        appBar: OnionBar("Участники", context),
        body: ListView.builder(
          itemCount: _competitorsFuture.length,
          itemBuilder: (context, index) {
            final competitor = _competitorsFuture[index].competitor;
            return Card(
              margin: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              child: ListTile(
                title: Text(competitor.fullName), // TODO style of the text
                subtitle: Text(
                  competitor.rank?.toString().split('.').last ??
                      "Разряд не указан",
                  style: TextStyle(color: Theme.of(context).primaryColorDark),
                ),
                leading: CircleAvatar(child: Text((index + 1).toString())),
              ),
            );
          },
        ),
      ),
    );
  }
}
