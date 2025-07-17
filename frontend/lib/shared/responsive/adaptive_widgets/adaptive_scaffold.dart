import 'package:flutter/material.dart';
import '../responsive.dart';

class AdaptiveScaffold extends StatelessWidget {
  final Widget body;
  final AppBar? appBar;
  final Widget? bottomNavigationBar;
  final Widget? drawer;
  final FloatingActionButton? floatingActionButton;
  final bool useScrollView;

  const AdaptiveScaffold({
    super.key,
    required this.body,
    this.appBar,
    this.bottomNavigationBar,
    this.drawer,
    this.floatingActionButton,
    this.useScrollView = true,
  });

  @override
  Widget build(BuildContext context) {
    final isDesktop = Responsive.isDesktop(context);

    return Scaffold(
      appBar: appBar,
      drawer: isDesktop ? null : drawer,
      bottomNavigationBar: isDesktop ? null : bottomNavigationBar,
      floatingActionButton: floatingActionButton,
      body: SafeArea(
        child: Row(
          children: [
            // Desktop Navigation Drawer
            if (isDesktop && drawer != null)
              SizedBox(
                width: 280,
                child: drawer!,
              ),

            // Main Content
            Expanded(
              child: useScrollView
                  ? SingleChildScrollView(
                      physics: const ClampingScrollPhysics(),
                      child: ConstrainedBox(
                        constraints: BoxConstraints(
                          minHeight: MediaQuery.of(context).size.height -
                              (appBar?.preferredSize.height ?? 0) -
                              MediaQuery.of(context).padding.top,
                        ),
                        child: body,
                      ),
                    )
                  : body,
            ),
          ],
        ),
      ),
    );
  }
}
